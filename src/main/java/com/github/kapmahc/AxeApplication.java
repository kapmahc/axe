package com.github.kapmahc;

import com.github.kapmahc.nut.resources.admin.LocaleController;
import com.github.kapmahc.nut.dao.LocaleDao;
import io.dropwizard.Application;
import io.dropwizard.db.DataSourceFactory;
import io.dropwizard.jdbi.DBIFactory;
import io.dropwizard.migrations.MigrationsBundle;
import io.dropwizard.setup.Bootstrap;
import io.dropwizard.setup.Environment;
import io.dropwizard.views.ViewBundle;
import org.glassfish.hk2.utilities.binding.AbstractBinder;
import org.skife.jdbi.v2.DBI;

public class AxeApplication extends Application<AxeConfiguration> {

    public static void main(final String[] args) throws Exception {
        new AxeApplication().run(args);
    }

    @Override
    public String getName() {
        return "axe";
    }

    @Override
    public void initialize(final Bootstrap<AxeConfiguration> bt) {
        bt.addBundle(new MigrationsBundle<AxeConfiguration>() {
            @Override
            public DataSourceFactory getDataSourceFactory(AxeConfiguration cfg) {
                return cfg.getDataSourceFactory();
            }
        });
        bt.addBundle(new ViewBundle<AxeConfiguration>() {

        });
    }

    @Override
    public void run(final AxeConfiguration cfg,
                    final Environment env) {
        final DBIFactory factory = new DBIFactory();
        final DBI jdbi = factory.build(env, cfg.getDataSourceFactory(), "postgresql");

        env.jersey().register(new AbstractBinder() {
            @Override
            protected void configure() {
                bind(cfg).to(AxeConfiguration.class);
                bind(env).to(Environment.class);
                bind(jdbi).to(DBI.class);
                bind(jdbi.onDemand(LocaleDao.class)).to(LocaleDao.class);
            }
        });

        env.jersey().register(LocaleController.class);
    }


}
