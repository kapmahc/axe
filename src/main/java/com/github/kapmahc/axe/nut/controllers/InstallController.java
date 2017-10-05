package com.github.kapmahc.axe.nut.controllers;

import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;

@Controller
@RequestMapping(value = "/install")
public class InstallController {
    @RequestMapping(method = RequestMethod.GET)
    public String getInstall(Model model) {
        return "nut/install";
    }

    @RequestMapping(method = RequestMethod.POST)
    public String postInstall(Model model, final RedirectAttributes redirectAttributes) {
        redirectAttributes.addFlashAttribute("notice", "aaa");
        return "redirect:/install";
    }
}
