import unittest

from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

from model import Post, Tag, Base

class TestOneToMany(unittest.TestCase):

    def setUp(self):
        self.engine = create_engine('sqlite://')
        DBSession = sessionmaker(bind=self.engine)
        self.session = DBSession()
        Base.metadata.create_all(self.engine)

    def tearDown(self):
        Base.metadata.drop_all(self.engine)

    def test_post_tag(self):
        t1 = Tag(name='iot')
        t2 = Tag(name='web')
        p1 = Post(content='This is a post contain iot and web tags.')
        p2 = Post(content='This is a post contain web.')
        p1.add_all_tags([t1, t2])
        p2.tags.append(t2)
        self.session.add_all([t1, t2,  p1, p2])
        self.session.commit()
        self.assertCountEqual(t1.posts, [p1])
        self.assertCountEqual(t2.posts, [p1, p2])
        self.assertCountEqual(p1.tags, [t1, t2])
        self.assertCountEqual(p2.tags, [t2])

if __name__ == '__main__':
    unittest.main()