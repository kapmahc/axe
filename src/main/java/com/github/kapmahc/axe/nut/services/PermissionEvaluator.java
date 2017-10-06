package com.github.kapmahc.axe.nut.services;

import com.github.kapmahc.axe.nut.models.User;
import com.github.kapmahc.axe.nut.repositories.UserRepository;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;

import javax.annotation.Resource;
import java.io.Serializable;

@Service
public class PermissionEvaluator implements org.springframework.security.access.PermissionEvaluator {
    @Override
    public boolean hasPermission(Authentication authentication, Object targetDomainObject, Object permission) {
        User user = userRepository.findByUid(authentication.getName());
        return user != null && userService.can(user, permission.toString(), targetDomainObject.toString(), null);
    }

    @Override
    public boolean hasPermission(Authentication authentication, Serializable targetId, String targetType, Object permission) {
        User user = userRepository.findByUid(authentication.getName());
        return user != null && userService.can(user, permission.toString(), targetType, (Long) targetId);
    }


    @Resource
    UserService userService;
    @Resource
    UserRepository userRepository;
}
