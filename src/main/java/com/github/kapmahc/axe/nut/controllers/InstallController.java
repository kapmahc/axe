package com.github.kapmahc.axe.nut.controllers;

import com.github.kapmahc.axe.nut.helper.FormHelper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Controller;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;

import javax.annotation.Resource;
import javax.validation.Valid;

@Controller
@RequestMapping(value = "/install")
public class InstallController {
    @GetMapping
    public String getInstall(InstallForm installForm) {
        return "nut/install";
    }

    @PostMapping
    public String postInstall(@Valid InstallForm installForm, BindingResult result, final RedirectAttributes attributes) {
        if (formHelper.check(result, attributes)) {

        }
        //
//        attributes.addFlashAttribute(NOTICE, "aaa");
        return "redirect:/install";

    }

    @Resource
    FormHelper formHelper;

    private final static Logger logger = LoggerFactory.getLogger(InstallController.class);
}
