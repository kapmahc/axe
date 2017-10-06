package com.github.kapmahc.axe.nut.controllers;

import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;

@Controller("nut.HomeController")
@PreAuthorize("permitAll")
public class HomeController {
    @GetMapping(value = "/")
    public String getInstall(Model model) {
        return "nut/home";
    }
}

