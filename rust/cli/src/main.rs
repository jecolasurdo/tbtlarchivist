#[global_allocator]
static ALLOC: snmalloc_rs::SnMalloc = snmalloc_rs::SnMalloc;

use std::io::Result;

fn main() -> Result<()> {
    println!("Hi. I don't do anything.");
    Ok(())
}