package com.github.kapmahc.axe.nut.helper;

import com.google.common.base.Joiner;
import com.google.common.base.Strings;
import org.springframework.stereotype.Component;
import org.springframework.validation.BindingResult;
import org.springframework.validation.FieldError;
import org.springframework.validation.ObjectError;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;

import javax.servlet.http.HttpServletRequest;
import java.util.ArrayList;
import java.util.List;

import static com.github.kapmahc.axe.Flash.ERROR;

@Component("nut.requestHelper")
public class RequestHelper {
    public String home(HttpServletRequest request){
        StringBuffer url = request.getRequestURL();
        String uri = request.getRequestURI();
        return url.substring(0, url.length()-uri.length());
    }
    public String clientIp(HttpServletRequest request) {
        String addr = request.getHeader("X-FORWARDED-FOR");
        return Strings.isNullOrEmpty(addr) ? request.getRemoteAddr() : addr;
    }

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
