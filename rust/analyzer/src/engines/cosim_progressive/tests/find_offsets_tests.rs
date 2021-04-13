// #![allow(clippy::needless_range_loop)]
// #![allow(clippy::too_many_lines)]

// use crate::engines::cosim_progressive::{new, Analyzer, Error, Settings};

// // additional cases:
// //  - candidate shorter than pass_one_sample_size returns error
// //  - candidate shorter than pass_two_sample_size returns error

// struct TestCase {
//     target: fn() -> Vec<i16>,
//     candidate: fn() -> Vec<i16>,
//     exp_result: fn() -> Result<Vec<i64>, Error>,
// }

// #[allow(clippy::needless_pass_by_value)]
// fn run_test_case(test_case: TestCase) {
//     let engine_settings = Settings {
//         target_sample_rate: 22_050,
//     };
//     let engine = new(engine_settings);
//     let target = (test_case.target)();
//     let candidate = (test_case.candidate)();
//     let exp_result = (test_case.exp_result)();

//     let act_result = engine.find_offsets(&candidate, &target);

//     match exp_result {
//         Err(_) => assert_eq!(act_result.is_err(), true, "expected error but no error",),
//         Ok(v) => assert_eq!(act_result.unwrap(), v),
//     }
// }

// /// candidate at head returns candidate
// #[test]
// fn single_candidate() {
//     run_test_case(TestCase {
//         target: || -> Vec<i16> {
//             let mut t = vec![0; 1024 * 10];
//             for i in 0..10 {
//                 t[i] = 1;
//             }
//             t
//         },
//         candidate: || -> Vec<i16> { vec![1; 10] },
//         exp_result: || -> Result<Vec<i64>, Error> { Ok(vec![0]) },
//     })
// }

// /// overlapping candidates attempts to return each.
// #[test]
// fn overlapping_candidates() {
//     run_test_case(TestCase {
//         target: || -> Vec<i16> {
//             let mut t = vec![0; 1024 * 10];
//             for i in 20..30 {
//                 t[i] = 1;
//             }
//             for i in 25..35 {
//                 t[i] = 1;
//             }
//             t
//         },
//         candidate: || -> Vec<i16> { vec![1; 10] },
//         exp_result: || -> Result<Vec<i64>, Error> { Ok(vec![907_029, 1_405_895]) },
//     })
// }

// /// immediately adjascent candidates returns each.
// #[test]
// fn adjascent_candidates() {
//     run_test_case(TestCase {
//         target: || -> Vec<i16> {
//             let mut t = vec![0; 1024 * 10];
//             for i in 20..30 {
//                 t[i] = 1;
//             }
//             for i in 31..41 {
//                 t[i] = 1;
//             }
//             for i in 42..52 {
//                 t[i] = 1;
//             }
//             t
//         },
//         candidate: || -> Vec<i16> { vec![1; 10] },
//         exp_result: || -> Result<Vec<i64>, Error> { Ok(vec![907_029, 1_405_895, 1_904_761]) },
//     })
// }

// /// multiple non-overlapping candidates returns all
// #[test]
// fn multiple_candidates() {
//     run_test_case(TestCase {
//         target: || -> Vec<i16> {
//             let mut t = vec![0; 1024 * 10];
//             for i in 20..30 {
//                 t[i] = 1;
//             }
//             for i in 40..50 {
//                 t[i] = 1;
//             }
//             for i in 60..70 {
//                 t[i] = 1;
//             }
//             t
//         },
//         candidate: || -> Vec<i16> { vec![1; 10] },
//         exp_result: || -> Result<Vec<i64>, Error> { Ok(vec![907_029, 1_814_058, 2_721_088]) },
//     })
// }

// /// candidate at tail returns candidate
// #[test]
// fn tail_candidate() {
//     run_test_case(TestCase {
//         target: || -> Vec<i16> {
//             let mut t = vec![0; 100];
//             for i in 89..99 {
//                 t[i] = 1;
//             }
//             t
//         },
//         candidate: || -> Vec<i16> { vec![1; 10] },
//         exp_result: || -> Result<Vec<i64>, Error> { Ok(vec![4_036_281]) },
//     })
// }

// /// candidate not present returns nothing
// #[test]
// fn candidate_not_present() {
//     run_test_case(TestCase {
//         target: || -> Vec<i16> { vec![0; 1024 * 10] },
//         candidate: || -> Vec<i16> { vec![1; 10] },
//         exp_result: || -> Result<Vec<i64>, Error> { Ok(vec![]) },
//     })
// }
