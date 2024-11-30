use anyhow::{Context, Result};
use clap::Parser;
use console::{style, Emoji};
use reqwest::{header, Url};
use std::{fs, path::{Path, PathBuf}};

/// Downloads the input for a particular day/year of advent of code.
/// 
/// It will create a folder named day_{d}_{year}, with a file named
/// input.txt.
/// 
/// If a template dir path is specified, it will copy all files within the folder
/// to day_{d}_{year}.
#[derive(Parser)]
struct Cli {
    /// The day of the problem
    #[arg(short = 'd', value_parser = clap::value_parser!(u8).range(1..=31))]
    day: u8,
    /// The year of the problem
    #[arg(short = 'y', value_parser = clap::value_parser!(u16).range(2015..=2024))]
    year: u16,
    /// Session cookie to be used for downloading the input. If this argument is provided, the cookie
    /// will be saved in a .session.lock file and future usages of this tool will look for that file.
    /// If the .session.lock file is missing, this argument MUST be provided
    #[arg(short = 's', verbatim_doc_comment)]
    seesion_cookie: Option<String>,
    /// Path to a template folder containing the boilerplate code that's being used for all problems
    #[arg(short = 't')]
    template_path: Option<PathBuf>,
}

static SESSION_LOCK_PATH: &str = ".session.lock";

static CLIP: Emoji<'_, '_> = Emoji("🔗  ", "");
static LOOKING_GLASS: Emoji<'_, '_> = Emoji("🔍  ", "");
static SMIRK: Emoji<'_, '_> = Emoji("😏  ", "");

fn get_problem_input(url: Url, session_cookie: String) -> Result<String> {
    let cookie_header = format!("session={}", session_cookie);
    let mut request_headers = header::HeaderMap::new();
    request_headers.insert(header::COOKIE, header::HeaderValue::from_str(&cookie_header)?);
    let client = reqwest::blocking::ClientBuilder::new().default_headers(request_headers).build().unwrap();
    
    let request = client.get(url.clone()).build()?;
    let result = client.execute(request)?;
    return result.text().with_context(|| format!("Could not fetch Advent of Code input at {}", url));
}

fn copy_dir_rec(src: &Path, dest: &Path) -> Result<()> {
    if !dest.exists() {
        fs::create_dir_all(&dest)?
    }
    for entry in fs::read_dir(&src).with_context(|| format!("Could not list dir: {}", src.display()))? {
        let entry = entry?;
        let file_type = entry.file_type()?;
        if file_type.is_dir() {
            copy_dir_rec(&entry.path(), &dest.join(entry.file_name()))?
        } else {
            fs::copy(entry.path(), dest.join(entry.file_name()))?;
        }
    }

    Ok(())
}

fn main() -> Result<()> {
    let args = Cli::parse();

    // ARGS VALIDATION //
    println!(
        "{} {}Validating args...",
        style("[1/2]").green().bold().dim(),
        LOOKING_GLASS
    );
    if let Some(session) = args.seesion_cookie {
        std::fs::write(SESSION_LOCK_PATH, session)
            .with_context(|| format!("Failed to write session cookie at {}", SESSION_LOCK_PATH))?;
    }

    let session_cookie = std::fs::read_to_string(SESSION_LOCK_PATH).with_context(|| 
        format!("There is no session cookie saved at {}, you must call the script with -s <<seesion-cookie>> to be able to download the prolem input", SESSION_LOCK_PATH))?;

    
    if let Some(template_path) = args.template_path.clone() {
        if !template_path.exists() {
            anyhow::bail!(format!("Template path pram is provided, yet there is not folder at {}", template_path.display()))
        }
        if !template_path.is_dir() {
            anyhow::bail!(format!("The template proided at {} is not a directory", template_path.display())) 
        }
    }

    // INPUT FETCHING //
    println!(
        "{} {}Downloading input and copying template...",
        style("[2/2]").green().bold().dim(),
        CLIP
    );
    let aoc_input_url = format!("https://adventofcode.com/{}/day/{}/input", args.year, args.day);
    let input = get_problem_input(Url::parse(&aoc_input_url)?, session_cookie)?;
    let problem_dir_path = PathBuf::from(format!("./day-{}-{}", args.day, args.year));
    fs::create_dir(&problem_dir_path).with_context(|| format!("Could not create folder at {}", problem_dir_path.display()))?;

    let input_path = problem_dir_path.join(PathBuf::from("input.txt"));
    fs::write(&input_path, input).with_context(|| format!("Could not write problem input at {}", input_path.display()))?;

    // TEMPLATE COPYING //
    if let Some(template_path) = args.template_path {
        copy_dir_rec(&template_path, &problem_dir_path).with_context(|| format!("Could not copy from {} to {}", template_path.display(), problem_dir_path.display()))?
    }

    println!(
        "{}{}", 
        style("You're good to go now, good luck with your problem!").green().bold().dim(), 
        SMIRK
    );
    Ok(())
}
