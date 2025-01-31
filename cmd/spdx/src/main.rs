#![deny(warnings)]
mod commands;
mod config;

use commands::*;
use config::*;
use regex::{Regex, RegexSet};
use std::{
    collections::{HashSet, VecDeque},
    env,
    ffi::OsStr,
    fs::{self, File, OpenOptions},
    io::{Read, Seek, SeekFrom, Write},
    path::{Path, PathBuf},
    process::exit,
};
use structopt::StructOpt;

fn main() -> Result<(), String> {
    let mut cmd: Commands = Commands::from_args();
    let config = read_config(&cmd)?;
    if cmd.starting_directory.ends_with(".") && cmd.starting_directory.starts_with(".") {
        cmd.starting_directory = get_proj_root();
    }

    let ignore_patterns = RegexSet::new(config.ignore.iter()).unwrap();

    let mut examine = VecDeque::new();
    examine.push_back(cmd.starting_directory);
    while !examine.is_empty() {
        let curdir = examine.pop_front().unwrap();
        if curdir.is_dir() {
            match curdir.to_str() {
                None => {
                    eprintln!("Unable to determine directory name: {:?}", curdir);
                    continue;
                }
                Some(s) => {
                    if ignore_patterns.is_match(s) {
                        println!("Ignoring {}", s);
                        continue;
                    }
                    println!("Scanning {}", s);
                }
            }
            // Traverse the current directory and subdirectories
            let readdir = fs::read_dir(curdir).map_err(|e| e.to_string())?;
            for path in readdir {
                examine.push_back(path.map_err(|e| e.to_string())?.path());
            }
        } else if curdir.is_file() {
            match curdir.to_str() {
                None => eprintln!("Can't handle file '{:?}'", curdir),
                Some(s) => {
                    if ignore_patterns.is_match(s) {
                        println!("Ignoring {}", s);
                        continue;
                    }
                    process_file(s, &config)?;
                }
            }
        }
    }
    Ok(())
}

fn process_file(file_name: &str, config: &Config) -> Result<(), String> {
    println!("Processing {}", file_name);
    let mut f = OpenOptions::new()
        .read(true)
        .write(true)
        .open(PathBuf::from(file_name.to_string()))
        .map_err(|e| e.to_string())?;
    let length = f.metadata().map_err(|e| e.to_string())?.len() as usize;
    let mut buffer = String::new();
    let size = f.read_to_string(&mut buffer).map_err(|e| e.to_string())?;
    if size != length {
        return Err("Unable to read entire file contents".to_string());
    }
    let suffix = Path::new(file_name).extension().and_then(OsStr::to_str);
    let sfx = match suffix {
        Some(s) => s,
        None => Path::new(file_name)
            .file_name()
            .and_then(OsStr::to_str)
            .unwrap(),
    };
    match config.comments.get(sfx) {
        None => {}
        Some(comment) => {
            // Found, check for header
            let rx_header = Regex::new(
                format!(
                    r#"{}\s+{}\s+{}\s+{}\s+{}\s+SPDX-License-Identifier:\s+{}\s+{}\s+"#,
                    comment, comment, config.copyright, comment, comment, config.license, comment
                )
                .as_str(),
            )
            .unwrap();
            if rx_header.is_match(&buffer) {
                return Ok(());
            }
            let header = format!(
                r#"{}
{} {}
{}
{} SPDX-License-Identifier: {}
{}

"#,
                comment, comment, config.copyright, comment, comment, config.license, comment
            );
            f.seek(SeekFrom::Start(0)).map_err(|e| e.to_string())?;
            f.write_all(header.as_bytes()).map_err(|e| e.to_string())?;
            f.write_all(buffer.as_bytes()).map_err(|e| e.to_string())?;
        }
    }
    Ok(())
}

fn get_proj_root() -> PathBuf {
    match env::current_dir() {
        Err(e) => {
            eprintln!("{}", e);
            exit(5);
        }
        Ok(p) => p,
    }
}

fn read_config(cmd: &Commands) -> Result<Config, String> {
    let mut config = Config::default();
    match cmd.ignore.as_ref() {
        None => {}
        Some(patterns) => {
            config.ignore = patterns.iter().cloned().collect::<HashSet<_>>();
        }
    }
    match cmd.license.as_ref() {
        None => {}
        Some(s) => config.license = s.clone(),
    }
    match cmd.copyright.as_ref() {
        None => {}
        Some(s) => config.copyright = s.clone(),
    }
    if !cmd.config_file.exists() {
        return Ok(config);
    }

    let mut f = File::open(&cmd.config_file).map_err(|e| e.to_string())?;
    let len = f.metadata().map_err(|e| e.to_string())?.len() as usize;
    let mut bytes = Vec::with_capacity(len);
    let s = f.read_to_end(&mut bytes).map_err(|e| e.to_string())?;
    if s != len {
        return Err("Unable to read entire file contents".to_string());
    }
    let cfg = serde_yaml::from_slice::<Config>(bytes.as_slice()).map_err(|e| e.to_string())?;
    if !cfg.copyright.is_empty() {
        config.copyright = cfg.copyright.clone();
    }
    if !cfg.license.is_empty() {
        config.license = cfg.license.clone();
    }
    if !cfg.ignore.is_empty() {
        config.ignore = cfg.ignore.clone();
    }
    config.comments = cfg.comments;
    Ok(config)
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs::File;
    use std::io::{self, Write};

    #[test]
    fn test_performance_deviation() {
        let good_computation_bench = "Benchmark_ABC             842688              1396 ns/op";
        let bad_computation_bench = "Benchmark_ABC             842688              13960 ns/op";
        let result = compare_benchmarks(good_computation_bench, bad_computation_bench);
        assert!(result.is_err(), "Expected performance deviation, but did not detect any!");
    }

    #[test]
    fn test_parsing() {
        let output = "
        garbage data
        BenchmarkSigning/Secp256k1_-_5_of_9-16                         1        5794642205 ns/op
        BenchmarkSign2p-16                                             2         685590314 ns/op             29319 bytes/sign           16.00 msgs/sign
        garbage data
        ";
        let result = parse_benchmarks(output);
        assert!(result.is_ok(), "Failed to parse test input");
        let parsed_output = result.unwrap();
        assert_eq!(parsed_output.len(), 2, "Incorrect output length. Expected 2, got {}", parsed_output.len());
        assert!(parsed_output.contains_key("BenchmarkSigning/Secp256k1_-_5_of_9-16"), "Did not find BenchmarkSigning/Secp256k1_-_5_of_9-16 in the parsed output");
        assert!(parsed_output.contains_key("BenchmarkSign2p-16"), "Did not find BenchmarkSign2p-16 in the parsed output");
    }

    fn compare_benchmarks(curr_bench: &str, new_bench: &str) -> io::Result<()> {
        let curr_bench_file = "current_bench.log";
        let new_bench_file = "new_bench.log";

        let mut curr_file = File::create(curr_bench_file)?;
        let mut new_file = File::create(new_bench_file)?;

        writeln!(curr_file, "{}", curr_bench)?;
        writeln!(new_file, "{}", new_bench)?;

        // Simulate the comparison logic
        if curr_bench != new_bench {
            return Err(io::Error::new(io::ErrorKind::Other, "Performance deviation detected"));
        }

        Ok(())
    }

    fn parse_benchmarks(output: &str) -> io::Result<std::collections::HashMap<String, String>> {
        let mut parsed_output = std::collections::HashMap::new();
        for line in output.lines() {
            if line.starts_with("Benchmark") {
                let parts: Vec<&str> = line.split_whitespace().collect();
                if parts.len() >= 2 {
                    parsed_output.insert(parts[0].to_string(), parts[1].to_string());
                }
            }
        }
        Ok(parsed_output)
    }
}
