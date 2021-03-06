use conv::prelude::*;

pub fn copy_slice<T>(dst: &mut [T], src: &[T])
where
    T: Copy,
{
    for (d, s) in dst.iter_mut().zip(src.iter()) {
        *d = *s;
    }
}

#[allow(clippy::as_conversions, clippy::cast_possible_truncation)]
pub fn scale_to_i16(v: f64) -> i16 {
    f64::round(v * f64::from(i16::MAX)) as i16
}

pub fn scale_from_i16(v: i16) -> f64 {
    f64::from(v) / f64::from(i16::MAX)
}

pub fn cosine_similarity(a: &[i16], b: &[i16]) -> f64 {
    sumdotproduct(a, b) / (a.sqrsum().sqrt() * b.sqrsum().sqrt())
}

pub fn rms(raw: &[i16]) -> f64 {
    (raw.sqrsum() / raw.len().value_as::<f64>().unwrap()).sqrt()
}

#[allow(clippy::as_conversions)]
pub fn index_to_nanoseconds(i: usize, sample_rate: usize) -> i64 {
    const NANOSECONDS: usize = 1_000_000_000;
    (i * NANOSECONDS / sample_rate) as i64
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
