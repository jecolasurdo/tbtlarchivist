#![allow(clippy::needless_range_loop)]
#![allow(clippy::too_many_lines)]

use crate::engines::cosine_similarity::{new, Analyzer, Error, Settings};

// additional cases:
//  - candidate shorter than pass_one_sample_size returns error
//  - candidate shorter than pass_two_sample_size returns error

struct TestCase {
    target: fn() -> Vec<i16>,
    candidate: fn() -> Vec<i16>,
    exp_result: fn() -> Result<Vec<i64>, Error>,
}

#[test]
fn single_candidate() {
    //  candidate at head returns candidate
    run(TestCase {
        target: || -> Vec<i16> {
            let mut t = vec![0; 1024 * 10];
            for i in 0..10 {
                t[i] = 1;
            }
            t
        },
        candidate: || -> Vec<i16> { vec![1; 10] },
        exp_result: || -> Result<Vec<i64>, Error> { Ok(vec![0]) },
    })
}

#[test]
fn overlapping_candidates() {
    // overlapping candidates attempts to return each.
    run(TestCase {
        target: || -> Vec<i16> {
            let mut t = vec![0; 1024 * 10];
            for i in 20..30 {
                t[i] = 1;
            }
            for i in 25..35 {
                t[i] = 1;
            }
            t
        },
        candidate: || -> Vec<i16> { vec![1; 10] },
        exp_result: || -> Result<Vec<i64>, Error> { Ok(vec![20, 31]) },
    })
}

#[test]
fn adjascent_candidates() {
    // immediately adjascent candidates returns each.
    run(TestCase {
        target: || -> Vec<i16> {
            let mut t = vec![0; 1024 * 10];
            for i in 20..30 {
                t[i] = 1;
            }
            for i in 31..41 {
                t[i] = 1;
            }
            for i in 42..52 {
                t[i] = 1;
            }
            t
        },
        candidate: || -> Vec<i16> { vec![1; 10] },
        exp_result: || -> Result<Vec<i64>, Error> { Ok(vec![20, 31, 42]) },
    })
}

#[test]
fn multiple_candidates() {
    run(TestCase {
        // multiple non-overlapping candidates returns all
        target: || -> Vec<i16> {
            let mut t = vec![0; 1024 * 10];
            for i in 20..30 {
                t[i] = 1;
            }
            for i in 40..50 {
                t[i] = 1;
            }
            for i in 60..70 {
                t[i] = 1;
            }
            t
        },
        candidate: || -> Vec<i16> { vec![1; 10] },
        exp_result: || -> Result<Vec<i64>, Error> { Ok(vec![20, 40, 60]) },
    })
}

#[test]
fn tail_candidate() {
    run(TestCase {
        // candidate at tail returns candidate
        target: || -> Vec<i16> {
            let mut t = vec![0; 100];
            for i in 89..99 {
                t[i] = 1;
            }
            t
        },
        candidate: || -> Vec<i16> { vec![1; 10] },
        exp_result: || -> Result<Vec<i64>, Error> { Ok(vec![89]) },
    })
}

#[test]
fn candidate_not_present() {
    run(TestCase {
        // candidate not present returns nothing
        target: || -> Vec<i16> { vec![0; 1024 * 10] },
        candidate: || -> Vec<i16> { vec![1; 10] },
        exp_result: || -> Result<Vec<i64>, Error> { Ok(vec![]) },
    })
}

#[allow(clippy::needless_pass_by_value)]
fn run(test_case: TestCase) {
    let engine_settings = Settings {
        pass_one_sample_size: 5,
        pass_one_threshold: 0.5,
        pass_two_sample_size: 5,
        pass_two_threshold: 0.7,
    };
    let engine = new(engine_settings);
    let target = (test_case.target)();
    let candidate = (test_case.candidate)();
    let exp_result = (test_case.exp_result)();

    let act_result = engine.find_offsets(&candidate, &target);

    match exp_result {
        Err(_) => assert_eq!(act_result.is_err(), true, "expected error but no error",),
        Ok(v) => assert_eq!(act_result.unwrap(), v),
    }
}
