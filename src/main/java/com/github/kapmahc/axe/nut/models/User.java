package com.github.kapmahc.axe.nut.models;

import com.google.common.io.BaseEncoding;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.DynamicUpdate;
import org.hibernate.annotations.UpdateTimestamp;

import javax.persistence.*;
import java.io.Serializable;
import java.io.UnsupportedEncodingException;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

@Entity
@Table(name = "users", indexes = {
        @Index(columnList = "name", name = "idx_users_name"),
        @Index(columnList = "email", unique = true, name = "idx_users_email"),
        @Index(columnList = "uid", unique = true, name = "idx_users_uid"),
        @Index(columnList = "providerType", name = "idx_users_provider_type"),
        @Index(columnList = "providerType,providerId", unique = true, name = "idx_users_provider"),
})
@DynamicUpdate
public class User implements Serializable {
    public enum Type {
        EMAIL, GOOGLE, FACEBOOK, WECHAT
    }

    // https://en.gravatar.com/site/implement/images/java/
    public void setGravatarLogo() throws NoSuchAlgorithmException, UnsupportedEncodingException {
        MessageDigest md = MessageDigest.getInstance("MD5");
        byte[] buf = md.digest(email.trim().toLowerCase().getBytes("CP1252"));
        logo = BaseEncoding.base16().lowerCase().encode(buf);
    }

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Long id;
    @Column(nullable = false, length = 64)
    private String name;
    @Column(nullable = false, updatable = false)
    private String email;
    @Column(nullable = false, length = 36)
    private String uid;
    private String password;
    @Column(nullable = false)
    private String providerId;
    @Column(nullable = false, length = 16)
    @Enumerated(EnumType.STRING)
    private Type providerType;
    private String logo;
    private long signInCount;
    @Column(length = 45)
    private String lastSignInIp;
    private Date lastSignInAt;
    @Column(length = 45)
    private String currentSignInIp;
    private Date currentSignInAt;
    private Date confirmedAt;
    private Date lockedAt;
    @Column(nullable = false)
    @UpdateTimestamp
    private Date updatedAt;
    @Column(nullable = false, updatable = false)
    @CreationTimestamp
    private Date createdAt;
    @OneToMany(mappedBy = "user")
    private List<Log> logs;
    @OneToMany(mappedBy = "user")
    private List<Attachment> attachments;
    @OneToMany(mappedBy = "user")
    private List<Policy> policies;


    public User() {
        logs = new ArrayList<>();
        attachments = new ArrayList<>();
        policies = new ArrayList<>();
    }

    public List<Policy> getPolicies() {
        return policies;
    }

    public void setPolicies(List<Policy> policies) {
        this.policies = policies;
    }

    public List<Attachment> getAttachments() {
        return attachments;
    }


    public void setAttachments(List<Attachment> attachments) {
        this.attachments = attachments;
    }

    public List<Log> getLogs() {
        return logs;
    }

    public void setLogs(List<Log> logs) {
        this.logs = logs;
    }

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

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

    public String getUid() {
        return uid;
    }

    public void setUid(String uid) {
        this.uid = uid;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getProviderId() {
        return providerId;
    }

    public void setProviderId(String providerId) {
        this.providerId = providerId;
    }

    public Type getProviderType() {
        return providerType;
    }

    public void setProviderType(Type providerType) {
        this.providerType = providerType;
    }

    public String getLogo() {
        return logo;
    }

    public void setLogo(String logo) {
        this.logo = logo;
    }

    public long getSignInCount() {
        return signInCount;
    }

    public void setSignInCount(long signInCount) {
        this.signInCount = signInCount;
    }

    public String getLastSignInIp() {
        return lastSignInIp;
    }

    public void setLastSignInIp(String lastSignInIp) {
        this.lastSignInIp = lastSignInIp;
    }

    public Date getLastSignInAt() {
        return lastSignInAt;
    }

    public void setLastSignInAt(Date lastSignInAt) {
        this.lastSignInAt = lastSignInAt;
    }

    public String getCurrentSignInIp() {
        return currentSignInIp;
    }

    public void setCurrentSignInIp(String currentSignInIp) {
        this.currentSignInIp = currentSignInIp;
    }

    public Date getCurrentSignInAt() {
        return currentSignInAt;
    }

    public void setCurrentSignInAt(Date currentSignInAt) {
        this.currentSignInAt = currentSignInAt;
    }

    public Date getConfirmedAt() {
        return confirmedAt;
    }

    public void setConfirmedAt(Date confirmedAt) {
        this.confirmedAt = confirmedAt;
    }

    public Date getLockedAt() {
        return lockedAt;
    }

    public void setLockedAt(Date lockedAt) {
        this.lockedAt = lockedAt;
    }

    public Date getUpdatedAt() {
        return updatedAt;
    }

    public void setUpdatedAt(Date updatedAt) {
        this.updatedAt = updatedAt;
    }

    public Date getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(Date createdAt) {
        this.createdAt = createdAt;
    }
}
