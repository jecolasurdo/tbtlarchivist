use crate::engines::cosine_similarity::{new, Analyzer, Settings};
use std::fs::File;
use std::io::Read;

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
