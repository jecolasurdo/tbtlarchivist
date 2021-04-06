use crate::engines::cosine_similarity::{cosine_similarity, i16_to_u32};

#[test]
fn cosine_similarity_happy_path() {
    let a: [i16; 8] = [2, 0, 1, 1, 0, 2, 1, 1];
    let b: [i16; 8] = [2, 1, 1, 0, 1, 1, 1, 1];
    let cs = cosine_similarity(&a, &b);
    assert_eq!(cs, 0.821_583_836_257_749_1)
}

#[test]
fn i16_to_u32_happy_path() {
    let test_cases = vec![
        (-32768_i16, 0_u32),
        (32767_i16, 65535_u32),
        (0_i16, 32768_u32),
        (10_i16, 32778_u32),
    ];
    for (input, exp_result) in test_cases {
        let act_result = i16_to_u32(input);
        assert_eq!(exp_result, act_result);
    }
}
