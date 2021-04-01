use crate::engines::cosine_similarity::cosine_similarity;

#[test]
fn happy_path() {
    let a: [i16; 8] = [2, 0, 1, 1, 0, 2, 1, 1];
    let b: [i16; 8] = [2, 1, 1, 0, 1, 1, 1, 1];
    let cs = cosine_similarity(&a, &b);
    assert_eq!(cs, 0.821_583_836_257_749_1)
}
