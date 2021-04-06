use analyzer::engines::cosine_similarity::{self, Settings};
use analyzer::engines::Analyzer;
use core::f64::consts::PI;
use criterion::{black_box, criterion_group, criterion_main, Criterion};
use std::fs::File;
use std::io::Read;

fn benchmark_mp3_to_raw(c: &mut Criterion) {
    let sample_path = String::from(
        "/Users/Joe/Documents/code/tbtlarchivist/rust/analyzer/benches/125ms_constant_192kbps_joint_stereo.mp3",
    );
    let mut file = File::open(sample_path).unwrap();
    let mut data = Vec::new();
    file.read_to_end(&mut data).unwrap();

    // values in engine_settings are irrevent to benchmark
    let engine_settings = Settings {
        pass_one_sample_size: 50,
        pass_one_threshold: 0.60,
        pass_two_sample_size: 500,
        pass_two_threshold: 0.8,
        // pass_one_sample_size: 9,
        // pass_one_threshold: 0.991,
        // pass_two_sample_size: 50,
        // pass_two_threshold: 0.99,
    };
    let engine = cosine_similarity::new(engine_settings);
    c.bench_function("mp3_to_raw", |b| {
        b.iter(|| engine.mp3_to_raw(black_box(&data)))
    });
}

#[allow(clippy::as_conversions)]
#[allow(clippy::cast_possible_truncation)]
fn benchmark_fingerprint(c: &mut Criterion) {
    let sample_rate_hz = 22_050;
    let freq_hz = 440;
    let duration_ms = 10;
    let samples_per_period = sample_rate_hz / freq_hz;
    let periods = sample_rate_hz * duration_ms / 1000 / samples_per_period;

    let mut signal: Vec<i16> = Vec::new();
    for _ in 1..periods {
        for i in 1..samples_per_period {
            signal.push((2.0 * PI / f64::from(i) * f64::from(i16::MAX)).floor() as i16);
        }
    }

    const NOT_APPLICABLE_USIZE: usize = 0;
    const NOT_APPLICABLE_F64: f64 = 0.0;
    let engine = cosine_similarity::new(Settings {
        pass_one_sample_size: NOT_APPLICABLE_USIZE,
        pass_one_threshold: NOT_APPLICABLE_F64,
        pass_two_sample_size: NOT_APPLICABLE_USIZE,
        pass_two_threshold: NOT_APPLICABLE_F64,
    });

    c.bench_function("fingerprint", |b| {
        b.iter(|| engine.fingerprint(black_box(&signal)))
    });
}

criterion_group!(benches, benchmark_mp3_to_raw, benchmark_fingerprint);
criterion_main!(benches);
