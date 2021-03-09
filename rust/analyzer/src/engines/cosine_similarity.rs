use crate::engines::Analyzer;
use anyhow::Result;
use minimp3::{Decoder, Error as MP3Error, Frame};
use rubato::{InterpolationParameters, InterpolationType, Resampler, SincFixedIn, WindowFunction};
use std::rc::Rc;

pub struct Engine {
    options: Settings,
}

pub struct Settings {
    pub pass_one_sample_density: usize,
    pub pass_one_sample_size: usize,
    pub pass_one_threshold: f64,
    pub pass_two_sample_size: usize,
    pub pass_two_threshold: f64,
}

#[allow(dead_code)]
pub fn new(options: Settings) -> Engine {
    Engine { options }
}

const TARGET_SAMPLE_RATE: f64 = 22_050.0;

impl Analyzer for Engine {
    /// Converts an mp3 to a 16bit mono raw audio vector.  The primary use case for this system is
    /// for podcasts, which are generally monaraul, so as a simple way of "converting" from stereo
    /// to mono this method just ignores the right track and returns audio from the left track.
    /// Each mp3 frame is decoded, has the right channel stripped, is resampled to
    /// `TARGET_SAMPLE_RATE`, and is then stitched with the following frame.
    fn mp3_to_raw(&self, mp3_bytes: &[u8]) -> Result<Vec<i16>> {
        let mut decoder = Decoder::new(mp3_bytes);
        loop {
            match decoder.next_frame() {
                Ok(Frame {
                    data,
                    sample_rate,
                    channels,
                    ..
                }) => {
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
                        mono_data.len(),
                        1,
                    );
                    
                    // Need to work through the following.  There are a few examples in the rubato
                    // codebase.  Will probably need to check those out.
                    // See: https://github.com/HEnquist/rubato/tree/master/examples`
                    todo!("resampler wants an Vec<f64>, but the raw audio is a Vec<i16>.")

                    let resampled_frame = resampler.process(&mono_data)?;
                }
                Err(MP3Error::Eof) => break,
                Err(e) => return Err(e.into()),
            }
        }

        // pub fn read_wav(filename: String) -> Result<Vec<i16>, std::io::Error> {
        //     let mut file = File::open(filename)?;
        //     let (_, bit_depth) = wav::read(&mut file)?;
        //     match bit_depth {
        //         wav::BitDepth::Sixteen(x) => Ok(x),
        //         _ => Err(std::io::Error::new(
        //             ErrorKind::Other,
        //             "Unexpected bit depth. Only support 16bit",
        //         )),
        //     }
        // }
        todo!()
    }

    fn phash(&self, _: &[i16]) -> Result<Vec<u8>> {
        todo!()
    }

    fn find_offsets(&self, candidate: &[i16], target: &[i16]) -> Result<Vec<i64>> {
        let windows = target.windows(candidate.len());
        let mut similarities: Vec<f64> = Vec::with_capacity(windows.len());
        let mut n = 0;
        let mut possibilities = Vec::new();
        for window in windows {
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
                possibilities.push(window);
            }
            similarities.push(cs);
        }

        let results = vec![];
        for w in possibilities {
            let cs = cosine_similarity(
                &w[..self.options.pass_two_sample_size],
                &candidate[..self.options.pass_two_sample_size],
            );
            if cs >= self.options.pass_two_threshold {
                // The POC only stored a vector of the final scores, but didn't store the position
                // of those scores. Will need to store that info in a map of some sort so the
                // offsets can be returned.
                todo!("should append result here")
            }
            similarities.push(cs);
        }

        Ok(results)
    }
}

pub fn to_monaural(data: &[i16], channels: usize) -> Result<Vec<i16>> {
    if !(1..=2).contains(&channels) {
        todo!("this needs to be returned as a proper error");
        panic!("unsupported number of channels encountered");
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
    Ok(mono)
}

pub fn cosine_similarity(a: &[i16], b: &[i16]) -> f64 {
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
