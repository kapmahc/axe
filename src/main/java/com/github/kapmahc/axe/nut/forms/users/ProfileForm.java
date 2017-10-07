package com.github.kapmahc.axe.nut.forms.users;

import javax.validation.constraints.NotNull;
import java.io.Serializable;

public class ProfileForm implements Serializable {
    private String email;
    @NotNull
    private String name;

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }
}
