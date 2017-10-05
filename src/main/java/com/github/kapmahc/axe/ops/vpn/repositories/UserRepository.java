package com.github.kapmahc.axe.ops.vpn.repositories;

import com.github.kapmahc.axe.ops.vpn.models.User;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("ops.vpn.userRepository")
public interface UserRepository extends CrudRepository<User, Long> {
}
