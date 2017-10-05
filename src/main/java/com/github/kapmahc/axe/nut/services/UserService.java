package com.github.kapmahc.axe.nut.services;

import com.github.kapmahc.axe.nut.repositories.UserRepository;
import org.springframework.stereotype.Service;

import javax.annotation.Resource;

@Service("nut.userService")
public class UserService {
    @Resource
    UserRepository userRepository;
}
