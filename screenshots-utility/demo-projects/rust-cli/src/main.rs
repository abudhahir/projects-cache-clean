use clap::Parser;

#[derive(Parser)]
#[command(name = "rust-cli")]
#[command(about = "A sample Rust CLI application")]
struct Cli {
    #[arg(short, long)]
    name: Option<String>,
}

fn main() {
    let cli = Cli::parse();
    
    match cli.name {
        Some(name) => println!("Hello, {}!", name),
        None => println!("Hello, World!"),
    }
}
