//! Provides business logic associated with cosine similarity analysis of audio samples.

use crate::engines::Analyzer;
use minimp3::{Decoder, Error as MP3Error, Frame};
use rubato::{InterpolationParameters, InterpolationType, Resampler, SincFixedIn, WindowFunction};
use thiserror::Error;

const TARGET_SAMPLE_RATE: i32 = 22_050;

/// Provides business logic associated with cosine similarity analysis of audio samples.
pub struct Engine {
    options: Settings,
}

/// Parameters that tune the behavior of an analysis.
pub struct Settings {
    /// Defines how many windows to evaluate.When making the first pass through
    /// the target audio.
    /// 0 => effectively produces no results
    /// 1 -> evaluates every window
    /// n -> evaluates every nth window
    pub pass_one_sample_density: usize,
    /// The number of contiguous samples compared between the candidate and
    /// target for each window in the initial "rough" pass.
    pub pass_one_sample_size: usize,
    /// The minimum cosine similarity value that must be met or exceeded to be
    /// considered a potential match.
    pub pass_one_threshold: f64,
    /// The number of contiguous samples compared between the candidate and
    /// target for each window in the second pass.
    pub pass_two_sample_size: usize,
    /// The minimum cosine similarity value that must be met or exceeded to be
    /// considered a match.
    pub pass_two_threshold: f64,
}

#[allow(dead_code)]
/// Constructs a new `Engine`
pub fn new(options: Settings) -> Engine {
    Engine { options }
}

impl Analyzer<Error> for Engine {
    /// Decodes an mp3 file to a 16bit mono raw audio vector.  The primary use case for this system
    /// is for podcasts, which are generally monaraul, so, as a simple way of "converting" from
    /// stereo to mono, this method just ignores one of the channels.  Each mp3 frame is decoded,
    /// has one channel stripped, is resampled to `TARGET_SAMPLE_RATE`, and is then stitched with
    /// the subsequent frame.
    #[inline]
    fn mp3_to_raw(&self, mp3_bytes: &[u8]) -> Result<Vec<i16>, Error> {
        let mut decoder = Decoder::new(mp3_bytes);
        let mut raw_data = vec![];
        let mut frames_buffer = vec![];
        let mut current_sample_rate: i32 = 0;
        loop {
            match decoder.next_frame() {
                Ok(Frame {
                    data,
                    sample_rate,
                    channels,
                    ..
                }) => {
                    let mut mono_data = to_monaural(&data, channels)?;
                    if sample_rate == current_sample_rate {
                        frames_buffer.append(&mut mono_data);
                    } else {
                        flush_buffer(current_sample_rate, &frames_buffer, &mut raw_data)?;
                        frames_buffer.clear();
                        current_sample_rate = sample_rate;
                    }
                }
                Err(MP3Error::Eof) => {
                    if !frames_buffer.is_empty() {
                        flush_buffer(current_sample_rate, &frames_buffer, &mut raw_data)?;
                    }
                    break;
                }
                Err(e) => return Err(Error(Box::new(ErrorKind::MiniMp3(e)))),
            }
        }

        Ok(raw_data.iter().map(|v| scale_to_i16(*v)).collect())
    }

    fn phash(&self, _: &[i16]) -> Result<Vec<u8>, Error> {
        Ok(vec![])
    }

    fn find_offsets(&self, candidate: &[i16], target: &[i16]) -> Result<Vec<i64>, Error> {
        let windows = target.windows(candidate.len());
        let mut n = 0;
        let mut possibilities = Vec::new();
        let mut offset_index = -1;
        for window in windows {
            offset_index += 1;
            n += 1;
            if n != self.options.pass_one_sample_density {
                continue;
            }
            n = 0;

            let cs = cosine_similarity(
                &window[..self.options.pass_one_sample_size],
                &candidate[..self.options.pass_one_sample_size],
            );
            if cs >= self.options.pass_one_threshold {
                possibilities.push((offset_index, window));
            }
        }

        let mut results = vec![];
        for (offset_index, window) in &possibilities {
            let cs = cosine_similarity(
                &window[..self.options.pass_two_sample_size],
                &candidate[..self.options.pass_two_sample_size],
            );
            if cs >= self.options.pass_two_threshold {
                results.push(*offset_index);
            }
        }

        Ok(results)
    }
}

fn flush_buffer(current_sample_rate: i32, src: &[f64], dst: &mut Vec<f64>) -> Result<(), Error> {
    let mut resampler = build_resampler(current_sample_rate, src.len());
    let mut cp = vec![];
    copy_slice(&mut cp, src);
    let fb = vec![cp; 1];
    let mut resampled_data = match resampler.process(&fb) {
        Ok(d) => d,
        Err(e) => return Err(Error(Box::new(ErrorKind::Resampler(e.to_string())))),
    };
    dst.append(&mut resampled_data[0]);
    Ok(())
}

fn build_resampler(sample_rate: i32, chunk_size: usize) -> impl rubato::Resampler<f64> {
    SincFixedIn::<f64>::new(
        f64::from(sample_rate) / f64::from(TARGET_SAMPLE_RATE),
        InterpolationParameters {
            sinc_len: 128,
            f_cutoff: 0.95,
            interpolation: InterpolationType::Nearest,
            oversampling_factor: 80,
            window: WindowFunction::BlackmanHarris2,
        },
        chunk_size,
        1,
    )
}

fn copy_slice<T>(dst: &mut [T], src: &[T])
where
    T: Copy,
{
    for (d, s) in dst.iter_mut().zip(src.iter()) {
        *d = *s;
    }
}

#[allow(clippy::as_conversions, clippy::cast_possible_truncation)]
fn scale_to_i16(v: f64) -> i16 {
    f64::round(v * f64::from(i16::MAX)) as i16
}

fn scale_from_i16(v: i16) -> f64 {
    f64::from(v) / f64::from(i16::MAX)
}

fn to_monaural(data: &[i16], channels: usize) -> Result<Vec<f64>, Error> {
    if !(1..=2).contains(&channels) {
        return Err(Error(Box::new(ErrorKind::InvalidChannelCount { channels })));
    }
    let mut mono: Vec<i16>;
    if channels == 2 {
        mono = Vec::with_capacity(data.len() / 2);
        let mut i = 0;
        while i < data.len() {
            mono.push(data[i]);
            i += 2;
        }
    } else {
        mono = data.to_vec();
    }
    Ok(mono.iter().map(|d| scale_from_i16(*d)).collect())
}

fn cosine_similarity(a: &[i16], b: &[i16]) -> f64 {
    sumdotproduct(a, b) / (a.sqrsum().sqrt() * b.sqrsum().sqrt())
}

fn sumdotproduct(a: &[i16], b: &[i16]) -> f64 {
    let mut sum = 0.0;
    for i in 0..a.len() {
        sum += f64::from(a[i]) * f64::from(b[i]);
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
            v += f64::from(*n) * f64::from(*n);
        }
        v
    }
}

/// A boxed error resulting from a problem running an engine.
#[derive(Error, Debug)]
#[error(transparent)]
pub struct Error(Box<ErrorKind>);

impl<E> From<E> for Error
where
    ErrorKind: From<E>,
{
    fn from(err: E) -> Self {
        Self(Box::new(ErrorKind::from(err)))
    }
}

/// Error variants associated with running an engine.
#[allow(missing_docs)]
#[derive(Error, Debug)]
pub enum ErrorKind {
    #[error("minimp3 error")]
    MiniMp3(#[from] minimp3::Error),

    #[error("resampler crate error: {0}")]
    Resampler(String),

    #[error("Inbound audio must have 1 or 2 channles, but contains {channels}.")]
    InvalidChannelCount { channels: usize },
}

#[cfg(test)]
mod tests {
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
            "/Users/Joe/Documents/code/tbtlarchivist/rust/analyzer/benches/drop_5000_samples.mp3",
        );
        let mut file = File::open(sample_path).unwrap();
        let mut data = Vec::new();
        file.read_to_end(&mut data).unwrap();

        // values in engine_settings are irrevent to this test
        let engine_settings = Settings {
            pass_one_sample_density: 1,
            pass_one_sample_size: 9,
            pass_one_threshold: 0.991,
            pass_two_sample_size: 50,
            pass_two_threshold: 0.99,
        };
        let engine = new(engine_settings);
        let result = engine.mp3_to_raw(&data).expect("should not panic");
        assert_eq!(1, result.len());
    }
}
