use crate::engines::cosim_progressive::index_to_nanoseconds;

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
