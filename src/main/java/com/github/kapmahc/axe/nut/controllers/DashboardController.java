package com.github.kapmahc.axe.nut.controllers;

import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

@Controller("auth.dashboardController")
@RequestMapping(value = "/dashboard")
public class DashboardController {
    @GetMapping
    public String index(Model model) {
        return "nut/dashboard";
    }
}
