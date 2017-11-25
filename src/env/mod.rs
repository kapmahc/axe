use clap::{App, Arg, ArgMatches, SubCommand};

pub mod app;
pub mod config;
pub mod errors;
pub mod utils;
pub mod router;

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
        .subcommand(
            SubCommand::with_name("database")
                .about("Database operations")
                .subcommand(SubCommand::with_name("migrate").about(
                    "Migrate the database to latest version.",
                ))
                .subcommand(SubCommand::with_name("rollback").about(
                    "Rollback the database to last version.",
                ))
                .subcommand(SubCommand::with_name("show").about(
                    "Show the database current version.",
                )),
        )
        .subcommand(
            SubCommand::with_name("cache")
                .about("Cache operations")
                .subcommand(SubCommand::with_name("list").about("List all cache keys."))
                .subcommand(SubCommand::with_name("clear").about(
                    "Clear all cache items.",
                )),
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
                    return app.database_migrate();
                }
                if let Some(_) = matches.subcommand_matches("rollback") {
                    return app.database_rollback();
                }
                if let Some(_) = matches.subcommand_matches("show") {
                    return app.database_show();
                }
                return Ok(());
            }
            if let Some(matches) = matches.subcommand_matches("cache") {
                if let Some(_) = matches.subcommand_matches("list") {
                    return app.cache_list();
                }
                if let Some(_) = matches.subcommand_matches("clear") {
                    return app.cache_clear();
                }
                return Ok(());
            }
            return app.start_server();
        }
    }
}
