package com.github.kapmahc.axe.nut.controllers;

import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;

@Controller
@RequestMapping(value = "/install")
public class InstallController {
    @RequestMapping(method = RequestMethod.GET)
    @PreAuthorize("permitAll")
    public String getInstall(Model model) {
        return "nut/install";
    }

    @PreAuthorize("permitAll")
    @RequestMapping(method = RequestMethod.POST)
    public void postInstall(Model model) {

    }
}
