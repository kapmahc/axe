package com.github.kapmahc.axe.nut.repositories;

import com.github.kapmahc.axe.nut.models.FriendLink;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("nut.friendLinkRepository")
public interface FriendLinkRepository extends CrudRepository<FriendLink, Long> {
}
