package com.github.kapmahc;

import io.dropwizard.Application;
import io.dropwizard.setup.Bootstrap;
import io.dropwizard.setup.Environment;

public class axeApplication extends Application<axeConfiguration> {

    public static void main(final String[] args) throws Exception {
        new axeApplication().run(args);
    }

    @Override
    public String getName() {
        return "axe";
    }

    @Override
    public void initialize(final Bootstrap<axeConfiguration> bootstrap) {
        // TODO: application initialization
    }

    @Override
    public void run(final axeConfiguration configuration,
                    final Environment environment) {
        // TODO: implement application
    }

}
