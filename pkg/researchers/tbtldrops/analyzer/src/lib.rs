use std::{fs::File, io::ErrorKind};
use thiserror::Error;

pub struct AnalyzerOptions {
    pub pass_one_sample_density: usize,
    pub pass_one_sample_size: usize,
    pub pass_one_threshold: f64,
    pub pass_two_sample_size: usize,
    pub pass_two_threshold: f64,
}

pub fn analyze(
    episode: Vec<i16>,
    drop: Vec<i16>,
    analyzer_options: AnalyzerOptions,
) -> Result<(), AnalyzerError> {
    let windows = episode.windows(drop.len());
    let mut similarities: Vec<f64> = Vec::with_capacity(windows.len());
    let mut i = 0;
    let mut n = 0;
    let mut possibilities = Vec::new();
    for window in windows {
        n += 1;
        if n != analyzer_options.pass_one_sample_density {
            continue;
        }
        n = 0;

        let cs = cosine_similarity(
            &window[..analyzer_options.pass_one_sample_size],
            &drop[..analyzer_options.pass_one_sample_size],
        );
        if i % 10_000_000 == 0 {
            println!("{}", i);
        }
        if cs >= analyzer_options.pass_one_threshold {
            possibilities.push(window);
        }
        similarities.push(cs);
        i += 1;
    }

    i = 0;
    println!("First pass found {} possibilities.", possibilities.len());
    for w in possibilities {
        let cs = cosine_similarity(
            &w[..analyzer_options.pass_two_sample_size],
            &drop[..analyzer_options.pass_two_sample_size],
        );
        if i % 1000 == 0 {
            println!("{}", i);
        }
        if cs >= analyzer_options.pass_two_threshold {
            println!("Found! {}", cs);
        }
        similarities.push(cs);
        i += 1;
    }

    Ok(())
}

pub fn cosine_similarity(a: &[i16], b: &[i16]) -> f64 {
    sumdotproduct(a, b) / (a.sqrsum().sqrt() * b.sqrsum().sqrt())
}

fn sumdotproduct(a: &[i16], b: &[i16]) -> f64 {
    let mut sum = 0.0;
    for i in 0..a.len() {
        sum += (a[i] as f64) * (b[i] as f64);
    }
    sum
}

trait SliceExt<T> {
    fn sqrsum(self) -> f64;
}

impl SliceExt<i16> for &[i16] {
    fn sqrsum(self) -> f64 {
        let mut v = 0.0;
        for n in self {
            v += (*n as f64) * (*n as f64);
        }
        v
    }
}

pub fn read_wav(filename: String) -> Result<Vec<i16>, std::io::Error> {
    let mut file = File::open(filename)?;
    let (_, bit_depth) = wav::read(&mut file)?;
    match bit_depth {
        wav::BitDepth::Sixteen(x) => Ok(x),
        _ => Err(std::io::Error::new(
            ErrorKind::Other,
            "Unexpected bit depth. Only support 16bit",
        )),
    }
}

#[derive(Error, Debug)]
pub enum AnalyzerError {
    #[error("channel receive error")]
    ChannelReceive(#[from] crossbeam_channel::RecvError),
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn cosine_sim_happy_path() {
        let a: [i16; 8] = [2, 0, 1, 1, 0, 2, 1, 1];
        let b: [i16; 8] = [2, 1, 1, 0, 1, 1, 1, 1];
        let cs = cosine_similarity(&a, &b);
        assert_eq!(cs, 0.8215838362577491)
    }
}
