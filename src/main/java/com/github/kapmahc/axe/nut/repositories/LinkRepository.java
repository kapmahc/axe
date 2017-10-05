package com.github.kapmahc.axe.nut.repositories;

import com.github.kapmahc.axe.nut.models.Link;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("nut.linkRepository")
public interface LinkRepository extends CrudRepository<Link, Long> {
}
