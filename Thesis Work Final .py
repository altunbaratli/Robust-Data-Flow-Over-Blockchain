#!/usr/bin/env python
# coding: utf-8

# In[34]:


from pyspark import SparkContext, SQLContext, SparkConf
from pyspark.sql import SparkSession
from pyspark.sql import Row
from pyspark.sql import functions as F
from pyspark.sql.functions import from_utc_timestamp
import os
import requests
import matplotlib.pyplot as plt


# In[35]:


# spark.stop()

conf = SparkConf().setAppName("PySpark App").setMaster("local").set("spark.jars", "/Users/altunbaratli/Desktop/Development/Robust-Data-Flow-Over-Blockchain").set('spark.driver.host','127.0.0.1')
sc = SparkContext(conf=conf)
sqlContext = SQLContext(sc)
spark = SparkSession     .builder     .appName("Python Spark SQL basic example").enableHiveSupport().getOrCreate()


# In[42]:


pload = {'username':'user2','channel':'channel1', 'smartcontract':'cc', 'args': {'sensorID':'1'}}
r = requests.post('http://169.51.200.103:30006/api/getHistory',json = pload)
print(r.text)


# In[43]:


json_rdd = sc.parallelize([r.text])
convertdf = spark.read.json(json_rdd)
convertdf.printSchema()


# In[45]:


ax = plt.gca()
newdf = convertdf.select("Value.sensorID", "Value.time", "Value.temp", "Value.stuckatfault", "Value.scfault", "Value.outlier", "Value.finaltemp")
flattened = newdf.dropna()
flattened_full = flattened.where("temp!=' '")
flattened_full2 = flattened_full.where("finaltemp!=' '")
df = flattened_full2.withColumn('temp',F.col('temp').cast('float'))
df2 = df.withColumn('finaltemp',F.col('finaltemp').cast('float'))
df2.toPandas().plot(kind='line', x="time", y="temp", figsize=(12,6), rot=50, ax=ax, label="Measurement SiD1")
# df2.toPandas().plot(kind='line', x="time", y="finaltemp", figsize=(12,6), rot=50, ax=ax, color="red", label="Result Temp.")
plt.show()


# In[ ]:




