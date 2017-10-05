package com.github.kapmahc.axe.nut.controllers;

import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;

@Controller
@PreAuthorize("permitAll")
public class HomeController {
    @RequestMapping(value = "/", method = RequestMethod.GET)
    public String getInstall(Model model) {
        return "nut/home";
    }
}

