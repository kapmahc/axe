package com.github.kapmahc.axe.nut.helper;

import com.google.common.base.Joiner;
import org.springframework.stereotype.Component;
import org.springframework.validation.BindingResult;
import org.springframework.validation.FieldError;
import org.springframework.validation.ObjectError;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;

import java.util.ArrayList;
import java.util.List;

import static com.github.kapmahc.axe.Flash.ERROR;

@Component
public class FormHelper {
    public boolean check(BindingResult result, final RedirectAttributes attributes) {
        if (!result.hasErrors()) {
            return true;
        }
        List<String> msg = new ArrayList<>();
        for (ObjectError e : result.getAllErrors()) {
            if (e instanceof FieldError) {
                FieldError oe = (FieldError) e;
                msg.add(oe.getField() + ": " + oe.getDefaultMessage());
            } else {
                msg.add(e.getDefaultMessage());
            }
        }
        attributes.addFlashAttribute(ERROR, Joiner.on("<br/>").join(msg));
        return false;
    }
}
