use clap::{App, Arg, ArgMatches, SubCommand};

pub mod app;
pub mod config;
pub mod errors;

pub const VERSION: &'static str = env!("CARGO_PKG_VERSION");

pub fn run() {
    let matches = App::new(env!("CARGO_PKG_DESCRIPTION"))
        .version(VERSION)
        .author(env!("CARGO_PKG_AUTHORS"))
        .about(env!("CARGO_PKG_HOMEPAGE"))
        .arg(
            Arg::with_name("config")
                .short("c")
                .long("config")
                .help("A custom config file")
                .default_value("config.toml")
                .takes_value(true),
        )
        // .arg(
        //     Arg::with_name("verbose")
        //         .short("v")
        //         .long("verbose")
        //         .help("Enable verbose output"),
        // )
        .subcommand(SubCommand::with_name("cache").about("Cache operations"))
        .subcommand(SubCommand::with_name("database").about("Database operations"))
        .subcommand(
            SubCommand::with_name("generate")
                .about("Generate files")
                .subcommand(SubCommand::with_name("config").about("Config file(toml)."))
                .subcommand(
                    SubCommand::with_name("locale")
                        .about("Locale file(yaml).")
                        .arg(
                            Arg::with_name("name")
                                .short("n")
                                .long("name")
                                .help("Locale name")
                                .required(true)
                                .takes_value(true),
                        ),
                )
                .subcommand(
                    SubCommand::with_name("migration")
                        .about("Database migration file(sql).")
                        .arg(
                            Arg::with_name("name")
                                .short("n")
                                .long("name")
                                .help("Migration name")
                                .required(true)
                                .takes_value(true),
                        ),
                )
                .subcommand(
                    SubCommand::with_name("nginx")
                        .about("Nginx config file.")
                        .arg(
                            Arg::with_name("https")
                                .short("s")
                                .long("https")
                                .help("HTTPS?"),
                        ),
                ),
        )
        .get_matches();

    match _main(&matches) {
        Ok(_) => {}
        Err(e) => error!("{}", e),
    }
}

fn _main<'a>(matches: &'a ArgMatches) -> errors::Result<()> {
    match matches.value_of("config") {
        None => Ok(()),
        Some(cfg) => {
            let app = app::App::new(cfg);
            if let Some(matches) = matches.subcommand_matches("generate") {
                if let Some(_) = matches.subcommand_matches("config") {
                    let it = config::Config::new();
                    return it.write(cfg.to_string());
                }
                if let Some(matches) = matches.subcommand_matches("locale") {
                    let name = format!("{}.yaml", matches.value_of("name").unwrap());
                    return app.generate_locale(&name);
                }
                if let Some(matches) = matches.subcommand_matches("migration") {
                    let name = matches.value_of("name").unwrap();
                    return app.generate_migration(&name);
                }
                if let Some(matches) = matches.subcommand_matches("nginx") {
                    let name = "nginx.conf";
                    let ssl = matches.is_present("https");
                    info!("generate file {} {}", name, ssl);
                    return Ok(());
                }

                return Ok(());
            }
            info!("start application server");
            Ok(())
        }
    }
}
