package com.github.kapmahc.axe.nut.models;

import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;

import javax.persistence.*;
import java.io.Serializable;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

@Entity
@Table(name = "roles", indexes = {
        @Index(columnList = "name,resourceType,resourceId", unique = true, name = "idx_roles_name_resource"),
        @Index(columnList = "name", name = "idx_roles_name"),
        @Index(columnList = "resourceType", name = "idx_roles_resource_type")
})
public class Role implements Serializable {
    public final static String ADMIN = "admin";
    public final static String ROOT = "root";
    public final static String MEMBER = "member";
    public final static String MANAGER = "manager";

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Long id;
    @Column(nullable = false)
    private String name;
    private Long resourceId;
    @Column(nullable = false)
    private String resourceType;
    @Column(nullable = false)
    @UpdateTimestamp
    private Date updatedAt;
    @Column(nullable = false, updatable = false)
    @CreationTimestamp
    private Date createdAt;
    @OneToMany(mappedBy = "role")
    private List<Policy> policies;


    public Role() {
        policies = new ArrayList<>();
    }

    public List<Policy> getPolicies() {
        return policies;
    }

    public void setPolicies(List<Policy> policies) {
        this.policies = policies;
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

    public Long getResourceId() {
        return resourceId;
    }

    public void setResourceId(Long resourceId) {
        this.resourceId = resourceId;
    }

    public String getResourceType() {
        return resourceType;
    }

    public void setResourceType(String resourceType) {
        this.resourceType = resourceType;
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
