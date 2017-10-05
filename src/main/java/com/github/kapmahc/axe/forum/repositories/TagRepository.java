package com.github.kapmahc.axe.forum.repositories;

import com.github.kapmahc.axe.forum.models.Tag;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("forum.tagRepository")
public interface TagRepository  extends CrudRepository<Tag, Long> {
}
