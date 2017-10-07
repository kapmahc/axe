package com.github.kapmahc.axe.nut.forms.users;

import javax.validation.constraints.NotNull;
import java.io.Serializable;

public class EmailForm implements Serializable {
    @NotNull
    private String email;

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }
}
