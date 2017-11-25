use clap::{App, Arg, ArgMatches, SubCommand};

pub mod app;
pub mod config;
pub mod errors;
pub mod utils;

pub fn run() {
    let matches = App::new(env!("CARGO_PKG_DESCRIPTION"))
        .version(env!("CARGO_PKG_VERSION"))
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
        .subcommand(SubCommand::with_name("cache").about("Cache operations"))
        .subcommand(SubCommand::with_name("database").about(
            "Database operations",
        ))
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
                        .arg(Arg::with_name("https").short("s").long("https").help(
                            "HTTPS?",
                        )),
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
            let app = app::App::new(cfg.to_string());
            if let Some(matches) = matches.subcommand_matches("generate") {
                if let Some(_) = matches.subcommand_matches("config") {
                    let c = config::Config::new();
                    return c.write(cfg.to_string());
                }
                if let Some(matches) = matches.subcommand_matches("locale") {
                    let name = format!("{}.yaml", matches.value_of("name").unwrap());
                    return app.generate_locale(name);
                }
                if let Some(matches) = matches.subcommand_matches("migration") {
                    let name = matches.value_of("name").unwrap();
                    return app.generate_migration(name.to_string());
                }
                if let Some(matches) = matches.subcommand_matches("nginx") {
                    let ssl = matches.is_present("https");
                    return app.generate_nginx_conf(ssl);
                }
                return Ok(());
            }

            if let Some(matches) = matches.subcommand_matches("database") {
                if let Some(_) = matches.subcommand_matches("migrate") {
                    let c = config::Config::new();
                    return c.write(cfg.to_string());
                }
                if let Some(matches) = matches.subcommand_matches("migrate") {
                    let name = format!("{}.yaml", matches.value_of("name").unwrap());
                    return app.generate_locale(name);
                }
                if let Some(matches) = matches.subcommand_matches("version") {
                    let name = matches.value_of("name").unwrap();
                    return app.generate_migration(name.to_string());
                }
                return Ok(());
            }
            return app.start_server();
        }
    }
}
