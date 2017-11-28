package com.github.kapmahc.nut.resources;

import io.dropwizard.views.View;

import javax.ws.rs.GET;
import javax.ws.rs.Path;

@Path("/users")
public class UsersController {
    @GET
    @Path("/sign-in")
    public View getSignIn(){
        return new View("sign-in.mustache"){} ;
    }

}
