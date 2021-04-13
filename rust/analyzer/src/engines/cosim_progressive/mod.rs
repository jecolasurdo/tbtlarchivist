//! Provides business logic associated with cosine similarity analysis of audio samples.

#[cfg(test)]
mod tests;

mod internals;

use crate::engines::cosim_progressive::internals::{
    copy_slice, index_to_nanoseconds, scale_from_i16, scale_to_i16,
};
use crate::engines::{Analyzer, Raw};
use minimp3::{Decoder, Error as MP3Error, Frame};
use rubato::{FftFixedIn, Resampler};
use std::convert::TryInto;
use thiserror::Error;

/// Provides business logic associated with cosine similarity analysis of audio samples.
pub struct Engine {
    options: Settings,
}

/// Parameters that tune the behavior of an analysis.
pub struct Settings {
    /// The sample rate to target when resampling inbound candidate and target
    /// audio. 22_050hz is a good starting point, as it retains audio up to
    /// 11khz.
    pub target_sample_rate: i32,
    /// The size of the window to use when calculating the peak RMS value for
    /// the candidate audio. 2756 samples is approx. 125ms at 22khz.
    pub rms_window_size: usize,
    /// The minimum cosine similarity value that must be met or exceeded to be
    /// considered a potential match.
    pub threshold: f64,
    /// The number of contiguous samples compared between the candidate and
    /// target for each window in the initial pass. This value is increased
    /// by a factor of 10 for each pass that the algorithm takes. For example,
    /// if `initial_sample_size` is set to 10, then the first pass compares
    /// 10 samples, second pass compares 100, third compares 1000, and so on
    /// until either the candidate is eliminated or the candidate is compared
    /// in full (submect to `max_sample_pct`).
    pub initial_sample_size: usize,
    /// The maximum percentage of the candidate audio to compare to the target.
    /// This value must be in the range [0,1)
    /// For example, if a candidate contains 10,000 datapoints, and
    /// `max_sample_pct` is set to 0.9, then up to only only 9,000 of the
    /// candidate's datapoints will be used to compare to the target.
    /// Specifically, the first and last 500 datapoints will be ignored when
    /// comparing the candidate to the target. This allows the algorithm to
    /// acknowledge that some percentage of candidate audio might be cropped.
    pub max_sample_pct: f64,
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
    /// has one channel stripped, is resampled to `Settings.target_sample_rate`, and is then stitched
    /// with the subsequent frame.
    #[inline]
    fn mp3_to_raw(&self, mp3_bytes: &[u8]) -> Result<Raw, Error> {
        let mut decoder = Decoder::new(mp3_bytes);
        let mut resampled_data = vec![];
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
                        resampled_data.append(&mut resample(
                            current_sample_rate,
                            self.options.target_sample_rate,
                            &frames_buffer,
                        )?);
                        frames_buffer.clear();
                    }
                    current_sample_rate = sample_rate;
                }
                Err(MP3Error::Eof) => {
                    if !frames_buffer.is_empty() {
                        resampled_data.append(&mut resample(
                            current_sample_rate,
                            self.options.target_sample_rate,
                            &frames_buffer,
                        )?);
                    }
                    break;
                }
                Err(e) => return Err(Error(Box::new(ErrorKind::MiniMp3(e)))),
            }
        }

        let data: Vec<i16> = resampled_data
            .iter()
            .map(|v| scale_to_i16(*v))
            .skip_while(|v| *v == 0)
            .collect();

        let duration_ns = index_to_nanoseconds(
            data.len(),
            self.options.target_sample_rate.try_into().unwrap(),
        );

        Ok(Raw { data, duration_ns })
    }

    /// Not implemented. Currently returns empty string.
    /// See commit f8d1df4: "Original fingerprint attempt" for original attempt
    /// which was prohibitively non-performant and never verified to work
    /// properly.
    fn fingerprint(&self, _raw: &[i16]) -> Result<String, Error> {
        Ok("".to_string())
    }

    /// Identifies positions within `target` where `candidate` is likely present.
    /// The resulting positions are expressed as nanoseconds.
    fn find_offsets(&self, _candidate: &[i16], _target: &[i16]) -> Result<Vec<i64>, Error> {
        unimplemented!();
    }
}

fn resample(
    current_sample_rate: i32,
    target_sample_rate: i32,
    buffer: &[f64],
) -> Result<Vec<f64>, Error> {
    let mut resampler = FftFixedIn::<f64>::new(
        current_sample_rate.try_into().unwrap(), // inbound sample rate
        target_sample_rate.try_into().unwrap(),  // desired sample rate
        buffer.len(),                            // frame size
        1024, // sub_chunks: this value is admittedly arbitrary. I'm not really sure how to rationalize it.
        1,    // number of channels
    );
    let mut cp = vec![0.0; buffer.len()];
    copy_slice(&mut cp, buffer);
    let fb = vec![cp; 1];
    match resampler.process(&fb) {
        Ok(d) => Ok(d[0].to_vec()),
        Err(e) => Err(Error(Box::new(ErrorKind::Resampler(e.to_string())))),
    }
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
