//! Provides business logic associated with cosine similarity analysis of audio samples.

use crate::engines::Analyzer;
use minimp3::{Decoder, Error as MP3Error, Frame};
use rubato::{InterpolationParameters, InterpolationType, Resampler, SincFixedIn, WindowFunction};
use thiserror::Error;

const TARGET_SAMPLE_RATE: f64 = 22_050.0;

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
    fn mp3_to_raw(&self, mp3_bytes: &[u8]) -> Result<Vec<i16>, Error> {
        let mut decoder = Decoder::new(mp3_bytes);
        let mut raw_data = vec![];
        loop {
            match decoder.next_frame() {
                Ok(Frame {
                    data,
                    sample_rate,
                    channels,
                    ..
                }) => {
                    // One thing that occurs to me is that we really only need to allocate
                    // a new resampler if/when the size of the frame changes or the frame's
                    // sample rate has changed,
                    // Look into https://github.com/bheisler/criterion.rs
                    todo!("the performance of this is really suspect.");
                    let params = InterpolationParameters {
                        sinc_len: 256,
                        f_cutoff: 0.95,
                        interpolation: InterpolationType::Nearest,
                        oversampling_factor: 160,
                        window: WindowFunction::BlackmanHarris2,
                    };
                    let mono_data = vec![to_monaural(&data, channels)?; 1];
                    let mut resampler = SincFixedIn::<f64>::new(
                        f64::from(sample_rate) / TARGET_SAMPLE_RATE,
                        params,
                        mono_data[0].len(),
                        1,
                    );
                    let mut resampled_data = match resampler.process(&mono_data) {
                        Ok(d) => d,
                        Err(e) => return Err(Error(Box::new(ErrorKind::Resampler(e.to_string())))),
                    };
                    raw_data.append(&mut resampled_data[0]);
                }
                Err(MP3Error::Eof) => break,
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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn cosine_sim_happy_path() {
        let a: [i16; 8] = [2, 0, 1, 1, 0, 2, 1, 1];
        let b: [i16; 8] = [2, 1, 1, 0, 1, 1, 1, 1];
        let cs = cosine_similarity(&a, &b);
        assert_eq!(cs, 0.821_583_836_257_749_1)
    }
}
