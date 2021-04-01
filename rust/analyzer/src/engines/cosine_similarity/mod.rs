//! Provides business logic associated with cosine similarity analysis of audio samples.

#[cfg(test)]
mod tests;

use crate::engines::Analyzer;
use minimp3::{Decoder, Error as MP3Error, Frame};
use rubato::{FftFixedIn, Resampler};
use std::convert::TryInto;
use std::ops::Neg;
use thiserror::Error;

const TARGET_SAMPLE_RATE: i32 = 22_050;

/// Provides business logic associated with cosine similarity analysis of audio samples.
pub struct Engine {
    options: Settings,
}

/// Parameters that tune the behavior of an analysis.
pub struct Settings {
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
        let mut current_sample_rate: i32 = -1;
        loop {
            // Resampling can leave undesireable artifacts at the frame
            // boundaries in the final audio data. Such artifacts can degrade
            // the effectiveness of the find_offsets process. To minimize frame
            // boundary artifacts, we append frame data to a buffer until the
            // sample rate for the next frame will differ from the buffer's
            // current sample rate.  The buffer is resampled and flushed only
            // when a sample rate change is detected (and when EOF is reached).
            match decoder.next_frame() {
                Ok(Frame {
                    data,
                    sample_rate,
                    channels,
                    ..
                }) => {
                    if current_sample_rate == -1 || sample_rate == current_sample_rate {
                        let mut mono_data = to_monaural(&data, channels)?;
                        frames_buffer.append(&mut mono_data);
                    } else {
                        raw_data.append(&mut resample(current_sample_rate, &frames_buffer)?);
                        frames_buffer.clear();
                    }
                    current_sample_rate = sample_rate;
                }
                Err(MP3Error::Eof) => {
                    if !frames_buffer.is_empty() {
                        raw_data.append(&mut resample(current_sample_rate, &frames_buffer)?);
                    }
                    break;
                }
                Err(e) => return Err(Error(Box::new(ErrorKind::MiniMp3(e)))),
            }
        }

        Ok(raw_data
            .iter()
            .map(|v| scale_to_i16(*v))
            .skip_while(|v| *v == 0)
            .collect())
    }

    fn phash(&self, _: &[i16]) -> Result<Vec<u8>, Error> {
        Ok(vec![])
    }

    #[allow(clippy::as_conversions)]
    fn find_offsets(&self, candidate: &[i16], target: &[i16]) -> Result<Vec<i64>, Error> {
        let windows = target.windows(candidate.len());
        let mut possibilities = Vec::new();
        let mut offset_index: i64 = -1;
        for window in windows {
            offset_index += 1;
            let cs = cosine_similarity(
                &window[..self.options.pass_one_sample_size],
                &candidate[..self.options.pass_one_sample_size],
            );
            if cs >= self.options.pass_one_threshold {
                possibilities.push((offset_index, window));
            }
        }

        let mut results = vec![];
        let mut cs_peak = f64::MIN;
        let mut i_peak = i64::neg(candidate.len() as i64);
        for (i_window, window) in &possibilities {
            // calculate score for current window
            let cs_window = cosine_similarity(
                &window[..self.options.pass_two_sample_size],
                &candidate[..self.options.pass_two_sample_size],
            );
            // if the current window's score doesn't meet the general threshold
            // continue to the next window.
            if cs_window < self.options.pass_two_threshold {
                continue;
            }

            // check to see if the current window index is outside the bounds
            // of a local peak.
            if *i_window > i_peak + (candidate.len() as i64) {
                // If we're here, we're ouside the bounds of a local peak
                // and are identifying a new local peak.
                // If a previous peak exists, push it to the result list.
                if cs_peak > f64::MIN {
                    results.push(i_peak);
                }
                // Set the local peak value and index to that of the current
                // window.
                cs_peak = cs_window;
                i_peak = *i_window;

                // move to the next window.
                continue;
            }

            // If we're here, we're within the bounds of a local peak.
            // Check to see if the current window value exceeds the current
            // peak.
            if cs_window > cs_peak {
                // Update the local peak value and index value.
                cs_peak = cs_window;
                i_peak = *i_window;
            }
        }
        // flush the last identified peak.
        if cs_peak > f64::MIN {
            results.push(i_peak);
        }

        Ok(results)
    }
}

fn resample(current_sample_rate: i32, buffer: &[f64]) -> Result<Vec<f64>, Error> {
    let mut resampler = build_resampler(current_sample_rate, buffer.len());
    let mut cp = vec![0.0; buffer.len()];
    copy_slice(&mut cp, buffer);
    let fb = vec![cp; 1];
    match resampler.process(&fb) {
        Ok(d) => Ok(d[0].to_vec()),
        Err(e) => Err(Error(Box::new(ErrorKind::Resampler(e.to_string())))),
    }
}

fn build_resampler(sample_rate: i32, chunk_size: usize) -> impl rubato::Resampler<f64> {
    FftFixedIn::<f64>::new(
        sample_rate.try_into().unwrap(),        // inbound sample rate
        TARGET_SAMPLE_RATE.try_into().unwrap(), // desired sample rate
        chunk_size,                             // frame size
        1024, // sub_chunks: this value is admittedly arbitrary. I'm not really sure how to rationalize it.
        1,    // number of channels
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
