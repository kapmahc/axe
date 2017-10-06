package com.github.kapmahc.axe.nut.forms.users;

import javax.validation.constraints.Size;
import java.io.Serializable;

public class ResetPasswordForm implements Serializable {
    @Size(min = 6, max = 32)
    private String password;
    private String passwordConfirmation;

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getPasswordConfirmation() {
        return passwordConfirmation;
    }

    public void setPasswordConfirmation(String passwordConfirmation) {
        this.passwordConfirmation = passwordConfirmation;
    }
}
