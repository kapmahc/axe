package com.github.kapmahc.axe.nut.forms.users;

import javax.validation.constraints.Size;
import java.io.Serializable;

public class SignInForm implements Serializable {
    @Size(min = 2, max = 255)
    private String email;
    @Size(min = 6, max = 32)
    private String password;

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }
}
