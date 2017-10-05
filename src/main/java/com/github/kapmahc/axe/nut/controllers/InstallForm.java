package com.github.kapmahc.axe.nut.controllers;


import javax.validation.constraints.NotNull;
import javax.validation.constraints.Size;
import java.io.Serializable;

public class InstallForm implements Serializable {
    @Size(min = 2, max = 32)
    private String subhead;
    @Size(min = 2, max = 255)
    private String title;
    @NotNull
    @Size(min = 2, max = 32)
    private String name;
    @NotNull
    @Size(min = 2, max = 32)
    private String email;
    @NotNull
    @Size(min = 6, max = 32)
    private String password;
    private String passwordConfirmation;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

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

    public String getPasswordConfirmation() {
        return passwordConfirmation;
    }

    public void setPasswordConfirmation(String passwordConfirmation) {
        this.passwordConfirmation = passwordConfirmation;
    }

    public String getSubhead() {
        return subhead;
    }

    public void setSubhead(String subhead) {
        this.subhead = subhead;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }
}
