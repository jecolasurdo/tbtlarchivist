use std::convert::TryFrom;

pub trait BytesExt<T> {
    fn to_frame(self) -> Vec<u8>;
}

impl BytesExt<u8> for &[u8] {
    /// Calculates the length of the slice as a 4 byte, big-endian i32 and
    /// returns a copy of the slice with the encoded length prepended to the
    /// original value.
    fn to_frame(self) -> Vec<u8> {
        let mut frame = i32::try_from(self.len()).unwrap().to_be_bytes().to_vec();
        frame.extend(&*self);
        frame
    }
}
