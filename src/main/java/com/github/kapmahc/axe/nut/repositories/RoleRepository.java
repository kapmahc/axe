package com.github.kapmahc.axe.nut.repositories;

import com.github.kapmahc.axe.nut.models.Role;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("nut.roleRepository")
public interface RoleRepository  extends CrudRepository<Role, Long> {
}
