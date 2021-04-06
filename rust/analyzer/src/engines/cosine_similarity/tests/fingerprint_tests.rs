use crate::engines::cosine_similarity::{new, Analyzer, Settings};
use std::f64::consts::PI;
#[test]
fn happy_path() {
    let raw = vec![1, 2, 3, 4, 5];
    let hash = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=".to_string();
    run(&raw, &hash);
}

#[test]
#[allow(clippy::as_conversions)]
#[allow(clippy::cast_possible_truncation)]
fn basic_waveform() {
    let sample_rate_hz = 22_050;
    let freq_hz = 440;
    let duration_sec = 1;
    let samples_per_period = sample_rate_hz / freq_hz;
    let periods = sample_rate_hz * duration_sec / samples_per_period;

    let mut signal: Vec<i16> = Vec::new();
    for _ in 1..periods {
        for i in 1..samples_per_period {
            signal.push((2.0 * PI / f64::from(i) * f64::from(i16::MAX)).floor() as i16);
        }
    }

    run(&signal, "foo");
}

fn run(raw: &[i16], exp_hash: &str) {
    const NOT_APPLICABLE_USIZE: usize = 0;
    const NOT_APPLICABLE_F64: f64 = 0.0;
    let engine = new(Settings {
        pass_one_sample_size: NOT_APPLICABLE_USIZE,
        pass_one_threshold: NOT_APPLICABLE_F64,
        pass_two_sample_size: NOT_APPLICABLE_USIZE,
        pass_two_threshold: NOT_APPLICABLE_F64,
    });

    let hash = engine
        .fingerprint(raw)
        .expect("test should not have errored");
    assert_eq!(hash, exp_hash);
}
