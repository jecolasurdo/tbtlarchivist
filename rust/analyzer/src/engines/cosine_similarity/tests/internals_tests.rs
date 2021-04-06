use crate::engines::cosine_similarity::{cosine_similarity, index_to_nanoseconds};

#[test]
fn cosine_similarity_happy_path() {
    let a: [i16; 8] = [2, 0, 1, 1, 0, 2, 1, 1];
    let b: [i16; 8] = [2, 1, 1, 0, 1, 1, 1, 1];
    let cs = cosine_similarity(&a, &b);
    assert_eq!(cs, 0.821_583_836_257_749_1)
}

#[test]
fn index_to_nanoseconds_happy_path() {
    let test_cases = vec![
        (0, 22_050, 0),
        (1, 22_050, 45_351),
        (21, 22_050, 952_380),
        (11_025, 22_050, 500_000_000),
        (22_050, 22_050, 1_000_000_000),
        (30_000, 22_050, 1_360_544_217),
        (79_380_000, 22_050, 3_600_000_000_000),
    ];
    for (index, sample_rate, exp_nanoseconds) in test_cases {
        let result = index_to_nanoseconds(index, sample_rate);
        assert_eq!(result, exp_nanoseconds);
    }
}
