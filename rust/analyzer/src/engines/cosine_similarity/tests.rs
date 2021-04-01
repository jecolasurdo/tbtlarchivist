#![allow(clippy::needless_range_loop)]
#![allow(clippy::too_many_lines)]

use super::*;
use std::fs::File;
use std::io::Read;

#[test]
fn cosine_sim_happy_path() {
    let a: [i16; 8] = [2, 0, 1, 1, 0, 2, 1, 1];
    let b: [i16; 8] = [2, 1, 1, 0, 1, 1, 1, 1];
    let cs = cosine_similarity(&a, &b);
    assert_eq!(cs, 0.821_583_836_257_749_1)
}
#[test]
fn mp3_to_raw_happy_path() {
    let sample_path = String::from(
            "/Users/Joe/Documents/code/tbtlarchivist/rust/analyzer/benches/125ms_constant_192kbps_joint_stereo.mp3",
        );
    let mut file = File::open(sample_path).unwrap();
    let mut data = Vec::new();
    file.read_to_end(&mut data).unwrap();

    // values in engine_settings are irrevent to this test
    let engine_settings = Settings {
        pass_one_sample_size: 9,
        pass_one_threshold: 0.991,
        pass_two_sample_size: 50,
        pass_two_threshold: 0.99,
    };
    let engine = new(engine_settings);
    let result = engine.mp3_to_raw(&data).expect("should not panic");
    // Ideally there should be 2756 samples in this result. However, the decoding and
    // resampling process prepends a bunch of zeros to the front and back of the outbound data.
    // Since the front of the audio is more important than the back of the audio, we trim any
    // zeros from the front and call it good.
    assert_eq!(3344, result.len());
}
// cases:
//  - candidate shorter than pass_one_sample_size returns error
//  - candidate shorter than pass_two_sample_size returns error
#[test]
fn find_offsets() {
    struct TestCase {
        name: String,
        target: fn() -> Vec<i16>,
        candidate: fn() -> Vec<i16>,
        exp_result: fn() -> Result<Vec<i64>, super::Error>,
    }

    let test_cases = vec![
        TestCase {
            //  single candidate present (not at head) returns candidate
            name: String::from("single candidate 1"),
            target: || -> Vec<i16> {
                let mut t = vec![0; 1024 * 10];
                for i in 20..30 {
                    t[i] = 1;
                }
                t
            },
            candidate: || -> Vec<i16> { vec![1; 10] },
            exp_result: || -> Result<Vec<i64>, super::Error> { Ok(vec![20]) },
        },
        TestCase {
            //  candidate at head returns candidate
            name: String::from("single candidate 2"),
            target: || -> Vec<i16> {
                let mut t = vec![0; 1024 * 10];
                for i in 0..10 {
                    t[i] = 1;
                }
                t
            },
            candidate: || -> Vec<i16> { vec![1; 10] },
            exp_result: || -> Result<Vec<i64>, super::Error> { Ok(vec![0]) },
        },
        TestCase {
            // overlapping candidates attempts to return each.
            name: String::from("overlapping candidates"),
            target: || -> Vec<i16> {
                let mut t = vec![0; 1024 * 10];
                for i in 20..30 {
                    t[i] = 1;
                }
                for i in 25..35 {
                    t[i] = 1;
                }
                t
            },
            candidate: || -> Vec<i16> { vec![1; 10] },
            exp_result: || -> Result<Vec<i64>, super::Error> { Ok(vec![20, 31]) },
        },
        TestCase {
            // immediately adjascent candidates returns each.
            name: String::from("immediately adjascent candidates"),
            target: || -> Vec<i16> {
                let mut t = vec![0; 1024 * 10];
                for i in 20..30 {
                    t[i] = 1;
                }
                for i in 31..41 {
                    t[i] = 1;
                }
                for i in 42..52 {
                    t[i] = 1;
                }
                t
            },
            candidate: || -> Vec<i16> { vec![1; 10] },
            exp_result: || -> Result<Vec<i64>, super::Error> { Ok(vec![20, 31, 42]) },
        },
        TestCase {
            // multiple non-overlapping candidates returns all
            name: String::from("multiple candidates"),
            target: || -> Vec<i16> {
                let mut t = vec![0; 1024 * 10];
                for i in 20..30 {
                    t[i] = 1;
                }
                for i in 40..50 {
                    t[i] = 1;
                }
                for i in 60..70 {
                    t[i] = 1;
                }
                t
            },
            candidate: || -> Vec<i16> { vec![1; 10] },
            exp_result: || -> Result<Vec<i64>, super::Error> { Ok(vec![20, 40, 60]) },
        },
        TestCase {
            // candidate at tail returns candidate
            name: String::from("tail candidate"),
            target: || -> Vec<i16> {
                let mut t = vec![0; 100];
                for i in 89..99 {
                    t[i] = 1;
                }
                t
            },
            candidate: || -> Vec<i16> { vec![1; 10] },
            exp_result: || -> Result<Vec<i64>, super::Error> { Ok(vec![89]) },
        },
        TestCase {
            // candidate not present returns nothing
            name: String::from("candidate not present"),
            target: || -> Vec<i16> { vec![0; 1024 * 10] },
            candidate: || -> Vec<i16> { vec![1; 10] },
            exp_result: || -> Result<Vec<i64>, super::Error> { Ok(vec![]) },
        },
    ];

    for test_case in test_cases {
        let engine_settings = Settings {
            pass_one_sample_size: 5,
            pass_one_threshold: 0.5,
            pass_two_sample_size: 5,
            pass_two_threshold: 0.7,
        };
        let engine = new(engine_settings);
        let target = (test_case.target)();
        let candidate = (test_case.candidate)();
        let exp_result = (test_case.exp_result)();

        let act_result = engine.find_offsets(&candidate, &target);

        match exp_result {
            Err(_) => assert_eq!(
                act_result.is_err(),
                true,
                "test_case '{}' expected error but no error",
                test_case.name
            ),
            Ok(v) => assert_eq!(act_result.unwrap(), v, "test case: '{}'", test_case.name),
        }
    }
}
#[ignore]
#[test]
fn mp3_to_raw_export() {
    let sample_path = String::from(
            "/Users/Joe/Documents/code/tbtlarchivist/rust/analyzer/benches/125ms_constant_192kbps_joint_stereo.mp3",
            // "/Users/Joe/Documents/code/tbtlarchivist/rust/audio/episodes/episode.mp3",
        );
    let mut file = File::open(sample_path).unwrap();
    let mut data = Vec::new();
    file.read_to_end(&mut data).unwrap();

    // values in engine_settings are irrevent to this test
    let engine_settings = Settings {
        pass_one_sample_size: 9,
        pass_one_threshold: 0.991,
        pass_two_sample_size: 50,
        pass_two_threshold: 0.99,
    };
    let engine = new(engine_settings);
    let raw = engine.mp3_to_raw(&data).expect("should not panic");

    let spec = hound::WavSpec {
        channels: 1,
        sample_rate: 22050,
        bits_per_sample: 16,
        sample_format: hound::SampleFormat::Int,
    };
    let mut writer = hound::WavWriter::create(
        "/Users/Joe/Documents/code/tbtlarchivist/rust/analyzer/benches/mp3_to_raw_export.wav",
        // "/Users/Joe/Documents/code/tbtlarchivist/rust/audio/episodes/episode_resampled_in_bulk.wav",
        spec,
    )
    .unwrap();
    for s in raw {
        writer.write_sample(s).unwrap();
    }
}
