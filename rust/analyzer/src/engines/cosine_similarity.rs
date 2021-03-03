
use crate::engines::EngineError;
use crate::engines::Analyzer;

pub struct Engine {
    options: Settings
}

pub struct Settings {
    pub pass_one_sample_density: usize,
    pub pass_one_sample_size: usize,
    pub pass_one_threshold: f64,
    pub pass_two_sample_size: usize,
    pub pass_two_threshold: f64,
}

pub fn new(options: Settings) -> Engine {
    Engine{
        options,
    } 
}

impl Analyzer for Engine {
    fn mp3_to_raw(&self, _: Vec<u8>) -> Result<Vec<i16>, EngineError> { 
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

    fn phash(&self, _: Vec<i16>) -> Result<Vec<u8>, EngineError> { 
        todo!() 
    }

    fn find_offsets(&self, candidate: Vec<i16>, target: Vec<i16>) -> Result<Vec<i64>, EngineError> { 
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
