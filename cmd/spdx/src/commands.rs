use std::path::PathBuf;
use structopt::StructOpt;

#[derive(Debug, StructOpt)]
pub struct Commands {
    #[structopt(short, long, parse(from_os_str), default_value = ".spdx.yml")]
    pub config_file: PathBuf,
    #[structopt(long)]
    pub copyright: Option<String>,
    #[structopt(short, long)]
    pub ignore: Option<Vec<String>>,
    #[structopt(long)]
    pub license: Option<String>,
    #[structopt(name = "DIR", parse(from_os_str), default_value = ".")]
    pub starting_directory: PathBuf,
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
