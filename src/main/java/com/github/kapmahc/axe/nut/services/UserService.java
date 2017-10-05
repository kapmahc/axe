package com.github.kapmahc.axe.nut.services;

import com.github.kapmahc.axe.nut.models.Role;
import com.github.kapmahc.axe.nut.repositories.LogRepository;
import com.github.kapmahc.axe.nut.repositories.PolicyRepository;
import com.github.kapmahc.axe.nut.repositories.RoleRepository;
import com.github.kapmahc.axe.nut.repositories.UserRepository;
import org.springframework.stereotype.Service;

import javax.annotation.Resource;
import java.time.ZoneId;
import java.util.Date;

@Service("nut.userService")
public class UserService {
    public void allow(long user, String role, int years) {
        Date begin = new Date();
        ZoneId zone = ZoneId.systemDefault();
        Date end = Date.from(begin.toInstant().atZone(zone).toLocalDateTime().plusYears(years).atZone(zone).toInstant());
        allow(user, role, null, null, new Date(), end);
    }


    public void allow(long user, String role, String resourceType, Long resourceId, Date begin, Date end) {

    }

    private Role getRole(String name, String resourceType, Long resourceId) {
        Role it = roleRepository.findByNameAndResourceTypeAndResourceId(name, resourceType, resourceId);
        if (it != null) {
            return it;
        }
        it = new Role();
        it.setName(name);
        it.setResourceType(resourceType);
        it.setResourceId(resourceId);
        roleRepository.save(it);
        return it;
    }

    @Resource
    UserRepository userRepository;
    @Resource
    RoleRepository roleRepository;
    @Resource
    PolicyRepository policyRepository;
    @Resource
    LogRepository logRepository;
}
