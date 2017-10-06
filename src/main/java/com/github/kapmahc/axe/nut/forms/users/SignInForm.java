package com.github.kapmahc.axe.nut.forms.users;

import javax.validation.constraints.NotNull;
import java.io.Serializable;

public class SignInForm implements Serializable {
    @NotNull
    private String email;
    @NotNull
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
